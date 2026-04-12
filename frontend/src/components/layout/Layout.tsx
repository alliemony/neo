import type { ReactNode } from 'react';
import { Header } from './Header';
import { Footer } from './Footer';

interface LayoutProps {
  children: ReactNode;
  sidebar?: ReactNode;
}

export function Layout({ children, sidebar }: LayoutProps) {
  return (
    <div className="min-h-screen bg-bg text-text-primary font-body">
      <Header />
      <div className="max-w-6xl mx-auto px-4 py-8 lg:flex lg:gap-8">
        <main className="flex-1 min-w-0">{children}</main>
        {sidebar}
      </div>
      <Footer />
    </div>
  );
}
