import type { ReactNode } from 'react';

interface SidebarProps {
  children: ReactNode;
}

export function Sidebar({ children }: SidebarProps) {
  return (
    <aside className="w-full lg:w-80 mt-8 lg:mt-0 space-y-6">
      {children}
    </aside>
  );
}
