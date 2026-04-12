import { BrowserRouter, Routes, Route } from 'react-router-dom';
import { Layout } from './components/layout/Layout';
import { Home } from './routes/Home';
import { PostView } from './routes/PostView';
import { TagFeed } from './routes/TagFeed';

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
        <Route path="/blog/:slug" element={<PostView />} />
        <Route path="/tag/:tag" element={<TagFeed />} />
        <Route path="*" element={<NotFound />} />
      </Routes>
    </BrowserRouter>
  );
}

export default App;
