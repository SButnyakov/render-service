import React from "react"
import SignupForm from "../components/SignupFormComponent/SignupFormComponent"

const styles = {
  backgroundImage: `url(${process.env.PUBLIC_URL + '/authBackground.png'})`,
  backgroundSize: "cover",
  height: "100vh",
  backgroundRepeat: "no-repeat",
}

export const RegisterPage = () => {
  return(
    <div style={styles}>
      <SignupForm/>
    </div>
  )
}

export default RegisterPage