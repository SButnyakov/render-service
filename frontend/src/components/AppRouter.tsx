import React from "react";
import { Navigate, Route, Routes } from "react-router-dom";
import LoginPage from "../pages/LoginPage";
import RegisterPage from "../pages/RegisterPage";

export const AppRouter = () => {
  return(
    <Routes>
      <Route path="/signin" element={<LoginPage/>}/>
      <Route path="/signup" element={<RegisterPage/>}/>
      <Route 
        path="*" 
        element={<Navigate to="/signin" replace/>}
      />
    </Routes>
  )
}

export default AppRouter;