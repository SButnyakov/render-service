import React from 'react'
import { BrowserRouter } from 'react-router-dom'
import AppRouter from './components/AppRouter'
import NavBar from './components/NavBar'
import { StoreProvider } from './hooks/useStore'

function App() {
  return (
    <StoreProvider>
      <BrowserRouter>
        <NavBar/>
        <AppRouter/>
      </BrowserRouter>
    </StoreProvider>
  )
}

export default App
