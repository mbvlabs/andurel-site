const COLORS = ['#f2ead8', '#e4dfd2', '#aaa393', '#8f8a7d'] as const

function mulberry32(seed: number) {
  return () => {
    let t = (seed += 0x6d2b79f5)
    t = Math.imul(t ^ (t >>> 15), t | 1)
    t ^= t + Math.imul(t ^ (t >>> 7), t | 61)
    return ((t ^ (t >>> 14)) >>> 0) / 4294967296
  }
}

function field(seed: number, count: number) {
  const rand = mulberry32(seed)
  return Array.from({ length: count }, (_, id) => {
    const faint = rand() > 0.18
    return {
      id,
      glyph: faint ? '·' : '.',
      color: COLORS[Math.floor(rand() * COLORS.length)],
      left: `${(rand() * 98 + 1).toFixed(2)}%`,
      top: `${(rand() * 96 + 1).toFixed(2)}%`,
      opacity: faint ? 0.1 + rand() * 0.16 : 0.22 + rand() * 0.18,
      size: faint ? 9 : 11,
    }
  })
}

const dots = field(0x2001, 96)

const galaxies = [
  {
    id: 'cluster-a',
    top: '11%',
    left: '76%',
    art: `    .   .  .
  .    .   .  .
 .  .     .   .
  .   .    .  .
    .  .   .`,
  },
  {
    id: 'cluster-b',
    top: '62%',
    left: '5%',
    art: `  .  .   .
 .   .  .  .
.  .     .  .
 .  .  .   .
  .   .  .`,
  },
  {
    id: 'cluster-c',
    top: '74%',
    left: '81%',
    art: `   .    .
 .   .    .  .
.  .    .   .
 .    .   .
   .    .`,
  },
]

export default function AsciiSky() {
  return (
    <div
      aria-hidden
      className="pointer-events-none fixed inset-0 overflow-hidden select-none font-mono"
    >
      {dots.map((dot) => (
        <span
          key={dot.id}
          className="absolute leading-none"
          style={{
            left: dot.left,
            top: dot.top,
            color: dot.color,
            opacity: dot.opacity,
            fontSize: dot.size,
          }}
        >
          {dot.glyph}
        </span>
      ))}
      {galaxies.map((galaxy) => (
        <pre
          key={galaxy.id}
          className="absolute text-[10px] leading-[1.2] tracking-[0.35em] text-[#aaa393]"
          style={{ top: galaxy.top, left: galaxy.left, opacity: 0.2 }}
        >
          {galaxy.art}
        </pre>
      ))}
    </div>
  )
}
