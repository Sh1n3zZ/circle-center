import { Routes, Route } from 'react-router-dom'
import Header from './components/layout/header'
import { Toaster } from './components/ui/sonner'
import './App.css'
import Home from './pages/Home'
import Reader from './pages/Reader'

const About = () => (
  <div className="p-4 text-center">
    <h1 className="text-3xl font-bold mb-4">About Page</h1>
    <p>This is a simple about page.</p>
  </div>
)

function App() {
  return (
      <div className="min-h-screen flex flex-col">
        <Header />

        <main className="flex-1">
          <Routes>
            <Route path="/" element={<Home />} />
            <Route path="/about" element={<About />} />
            <Route path="/reader" element={<Reader />} />
          </Routes>
        </main>

        <Toaster richColors position="top-right" />
      </div>
  )
}

export default App
