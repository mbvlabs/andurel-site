const searchDialog = document.querySelector("#docs-search")
const searchInput = document.querySelector("#docs-search-input")
const searchResults = document.querySelector("#docs-search-results")
const navDialog = document.querySelector("#docs-nav")
let searchIndex

const openSearch = async () => {
  searchDialog?.showModal()
  searchInput?.focus()
  if (!searchIndex) {
    searchIndex = fetch("/docs/search.json").then((response) => response.json())
  }
}

document.querySelectorAll("[data-docs-search-open]").forEach((button) => {
  button.addEventListener("click", openSearch)
})

document.querySelector("[data-docs-nav-open]")?.addEventListener("click", () => navDialog?.showModal())

document.querySelectorAll("[data-dialog-close]").forEach((button) => {
  button.addEventListener("click", () => button.closest("dialog")?.close())
})

document.addEventListener("keydown", (event) => {
  if ((event.metaKey || event.ctrlKey) && event.key.toLowerCase() === "k") {
    event.preventDefault()
    openSearch()
  }
})

searchInput?.addEventListener("input", async () => {
  const query = searchInput.value.trim().toLowerCase()
  if (!query) {
    searchResults.innerHTML = '<p class="px-3 py-8 text-center text-sm text-[#8f8a7d]">Start typing to search the documentation.</p>'
    return
  }

  const matches = (await searchIndex)
    .filter((item) => `${item.title} ${item.section} ${item.description} ${item.text}`.toLowerCase().includes(query))
    .slice(0, 10)

  searchResults.replaceChildren()
  if (!matches.length) {
    const empty = document.createElement("p")
    empty.className = "px-3 py-8 text-center text-sm text-[#8f8a7d]"
    empty.textContent = "No documentation found."
    searchResults.append(empty)
    return
  }

  matches.forEach((item) => {
    const link = document.createElement("a")
    link.href = item.url
    link.className = "block border border-transparent px-3 py-3 transition hover:border-[#2f3a37] hover:bg-[#171c1b] focus:border-[#8df7a4] focus:outline-none"

    const section = document.createElement("span")
    section.className = "block font-mono text-[0.65rem] uppercase tracking-wider text-[#8df7a4]"
    section.textContent = item.section

    const title = document.createElement("strong")
    title.className = "mt-1 block text-sm text-[#f2ead8]"
    title.textContent = item.title

    const description = document.createElement("span")
    description.className = "mt-1 block text-xs leading-5 text-[#8f8a7d]"
    description.textContent = item.description

    link.append(section, title, description)
    searchResults.append(link)
  })
})

searchResults?.addEventListener("keydown", (event) => {
  if (event.key !== "ArrowDown" && event.key !== "ArrowUp") return
  const links = [...searchResults.querySelectorAll("a")]
  const current = links.indexOf(document.activeElement)
  const next = event.key === "ArrowDown"
    ? Math.min(current + 1, links.length - 1)
    : Math.max(current - 1, 0)
  links[next]?.focus()
  event.preventDefault()
})

document.querySelectorAll(".docs-content pre").forEach((pre) => {
  const button = document.createElement("button")
  button.type = "button"
  button.className = "copy-code"
  button.textContent = "Copy"
  button.addEventListener("click", async () => {
    const code = pre.querySelector("code")?.cloneNode(true)
    code?.querySelectorAll('[style*="user-select:none"]').forEach((lineNumber) => lineNumber.remove())
    await navigator.clipboard.writeText(code?.textContent || pre.textContent)
    button.textContent = "Copied"
    setTimeout(() => { button.textContent = "Copy" }, 1500)
  })
  pre.append(button)
})

document.querySelectorAll(".docs-content h2[id], .docs-content h3[id]").forEach((heading) => {
  const link = document.createElement("a")
  link.href = `#${heading.id}`
  link.className = "ml-2 no-underline opacity-0 transition group-hover:opacity-100 focus:opacity-100"
  link.setAttribute("aria-label", `Link to ${heading.textContent}`)
  link.textContent = "#"
  heading.classList.add("group")
  heading.append(link)
})
