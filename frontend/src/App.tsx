import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/layout/Layout';

function Home() {
  return (
    <Layout>
      <h1 className="font-heading text-3xl font-bold mb-4">neo</h1>
      <p className="text-text-secondary">personal web garden</p>
    </Layout>
  );
}

function NotFound() {
  return (
    <Layout>
      <div className="text-center py-20">
        <h1 className="font-heading text-6xl mb-4">404</h1>
        <p className="text-text-secondary">Page not found.</p>
      </div>
    </Layout>
  );
}

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/" element={<Home />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
