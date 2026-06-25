import { BrowserRouter, Routes, Route, Link } from "react-router-dom";
import { Layout } from "./components/layout/Layout";
import { AuthProvider } from "./hooks/useAuth";
import { Home } from "./routes/Home";
import { Blog } from "./routes/Blog";
import { PostView } from "./routes/PostView";
import { About } from "./routes/About";
import { Recs } from "./routes/Recs";
import { Widgets } from "./routes/Widgets";
import { PageView } from "./routes/PageView";
import { AdminLogin } from "./routes/AdminLogin";
import { AdminDashboard } from "./routes/AdminDashboard";
import { PostEditor } from "./routes/PostEditor";
import { PageEditor } from "./routes/PageEditor";
import { ProtectedRoute } from "./routes/ProtectedRoute";

function NotFound() {
  return (
    <Layout>
      <div className="text-center py-20">
        <h1 className="text-6xl font-bold mb-4 text-text-primary">404</h1>
        <p className="text-text-secondary mb-4">The page you're looking for doesn't exist.</p>
        <Link to="/" className="text-accent hover:text-accent-alt">← Back to home</Link>
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
          <Route path="/blog" element={<Blog />} />
          <Route path="/blog/:slug" element={<PostView />} />
          <Route path="/about" element={<About />} />
          <Route path="/recs" element={<Recs />} />
          <Route path="/widgets" element={<Widgets />} />
          <Route path="/widgets/:id" element={<Widgets />} />
          <Route path="/page/:slug" element={<PageView />} />
          <Route path="/admin/login" element={<AdminLogin />} />
          <Route
            path="/admin"
            element={
              <ProtectedRoute>
                <AdminDashboard />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/posts/new"
            element={
              <ProtectedRoute>
                <PostEditor />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/posts/:slug/edit"
            element={
              <ProtectedRoute>
                <PostEditor />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/pages/new"
            element={
              <ProtectedRoute>
                <PageEditor />
              </ProtectedRoute>
            }
          />
          <Route
            path="/admin/pages/:slug/edit"
            element={
              <ProtectedRoute>
                <PageEditor />
              </ProtectedRoute>
            }
          />
          <Route path="*" element={<NotFound />} />
        </Routes>
      </BrowserRouter>
    </AuthProvider>
  );
}

export default App;
