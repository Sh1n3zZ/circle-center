import { Routes, Route } from 'react-router-dom'
import Header from './components/layout/header'
import { Toaster } from './components/ui/sonner'
import './App.css'
import Home from './pages/Home'
import Login from './pages/Login'
import Profile from './pages/profile/Profile'
import ProjectsPage from './pages/manager/Projects'
import ManagerProjectDetail from './pages/manager/Detail'

function App() {
  return (
      <div className="min-h-screen flex flex-col">
        <Header />

        <main className="flex-1">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/login" element={<Login />} />
            <Route path="/profile" element={<Profile />} />
            <Route path="/manager/projects" element={<ProjectsPage />} />
            <Route path="/manager/projects/:id" element={<ManagerProjectDetail />} />
          </Routes>
        </main>

        <Toaster richColors position="top-right" />
      </div>
  )
}

export default App
