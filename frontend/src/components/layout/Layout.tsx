import type { ReactNode } from 'react';
import { Header } from './Header';
import { Footer } from './Footer';

interface LayoutProps {
  children: ReactNode;
}

export function Layout({ children }: LayoutProps) {
  return (
    <div className="min-h-screen bg-bg text-text-secondary font-body">
      <Header />
      <main className="max-w-[700px] mx-auto px-6 py-5 pb-20">
        {children}
      </main>
      <Footer />
    </div>
  );
}
