import Header from "@/components/Header";
import DocsSidebar from "@/components/DocsSidebar";
import { SidebarProvider } from "@/components/SidebarContext";

export default function DocsLayout({
  children,
}: {
  children: React.ReactNode;
}) {
  return (
    <SidebarProvider>
      <Header />
      <div className="max-w-[1500px]  flex">
        <DocsSidebar />
        <main className="flex-1 min-w-0 py-8 px-6 lg:px-12">{children}</main>
      </div>
    </SidebarProvider>
  );
}
