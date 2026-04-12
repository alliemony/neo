import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import { Layout } from "./components/layout/Layout";
import { AuthProvider } from "./hooks/useAuth";
import { Home } from "./routes/Home";
import { PostView } from "./routes/PostView";
import { TagFeed } from "./routes/TagFeed";
import { PageView } from "./routes/PageView";
import { WidgetView } from "./routes/WidgetView";
import { AdminLogin } from "./routes/AdminLogin";
import { AdminDashboard } from "./routes/AdminDashboard";
import { PostEditor } from "./routes/PostEditor";
import { PageEditor } from "./routes/PageEditor";

function NotFound() {
  return (
    <Layout>
      <div className="text-center py-20">
        <h1 className="font-heading text-6xl mb-4">404</h1>
        <p className="text-text-secondary mb-4">
          The page you're looking for doesn't exist.
        </p>
        <Link to="/" className="text-accent hover:underline font-heading">
          ← Back to home
        </Link>
      </div>
    </Layout>
  );
}

function App() {
  return (
    <AuthProvider>
      <BrowserRouter>
        <Routes>
          <Route path="/" element={<Home />} />
          <Route path="/blog/:slug" element={<PostView />} />
          <Route path="/tag/:tag" element={<TagFeed />} />
          <Route path="/page/:slug" element={<PageView />} />
          <Route path="/widgets/:id" element={<WidgetView />} />
          <Route path="/admin/login" element={<AdminLogin />} />
          <Route path="/admin" element={<AdminDashboard />} />
          <Route path="/admin/posts/new" element={<PostEditor />} />
          <Route path="/admin/posts/:slug/edit" element={<PostEditor />} />
          <Route path="/admin/pages/new" element={<PageEditor />} />
          <Route path="/admin/pages/:slug/edit" element={<PageEditor />} />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
