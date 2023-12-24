import React, { useContext } from "react"
import { Navigate, Route, Routes } from "react-router-dom"
import AuthPage from "../pages/AuthPage"
import RegisterPage from "../pages/RegisterPage"
import UploadPage from "../pages/UploadPage"
// import RegisterPage from "../pages/RegisterPage"
// import UploadPage from "../pages/UploadPage"

export const AppRouter = () => {

  return(
    <Routes>
      <Route path="/signin" element={<AuthPage/>}/>
      <Route path="/signup" element={<RegisterPage/>}/>
      <Route path="/upload" element={<UploadPage/>}/>
      <Route 
        path="*" 
        element={<Navigate to="/signin" replace/>}
      />
    </Routes>
  )
}

export default AppRouter