import DashboardLayout from '@/Layouts/DashboardLayout'

export default function Home() {
  return (
    <DashboardLayout title="Overview">
      <p className="max-w-xl text-sm leading-6 text-muted-foreground">
        Admin console home. Additional sections will land in the sidebar as they are added.
      </p>
    </DashboardLayout>
  )
}
