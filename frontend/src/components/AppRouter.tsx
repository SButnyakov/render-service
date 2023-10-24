import React from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import LoginPage from "../pages/LoginPage";
import RegisterPage from "../pages/RegisterPage";
import UploadPage from "../pages/UploadPage";

export const AppRouter = () => {
  return(
    <Routes>
      <Route path="/signin" element={<LoginPage/>}/>
      <Route path="/signup" element={<RegisterPage/>}/>
      <Route path="/upload" element={<UploadPage/>}/>
      <Route 
        path="*" 
        element={<Navigate to="/signin" replace/>}
      />
    </Routes>
  )
}

export default AppRouter;