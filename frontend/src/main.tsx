import React from 'react'
import ReactDOM from 'react-dom/client'
import { createBrowserRouter, RouterProvider } from 'react-router-dom'
import App from './App'
import Login from './pages/Login'
import Register from './pages/Register'
import Chat from './pages/Chat'
import './index.css'


const router = createBrowserRouter([
{ path: '/', element: <App/> },
{ path: '/login', element: <Login/> },
{ path: '/register', element: <Register/> },
{ path: '/chat', element: <Chat/> },
])


ReactDOM.createRoot(document.getElementById('root')!).render(
<React.StrictMode>
<RouterProvider router={router} />
</React.StrictMode>,
)