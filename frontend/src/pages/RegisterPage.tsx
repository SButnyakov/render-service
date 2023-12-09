import React from "react"
import SignupForm from "../components/SignupFormComponent/SignupFormComponent"

const backgroundImagePath = '../staticObjects/images/auth_background.png'

export const RegisterPage = () => {
  return(
    <div style={{
      backgroundImage: `url(/auth_background.png)`,
      backgroundRepeat: "no-repeat",
      backgroundSize: "cover",
      width: '300px',
      height: '300px'
      }}>
      <SignupForm/>
    </div>
  )
}

export default RegisterPage