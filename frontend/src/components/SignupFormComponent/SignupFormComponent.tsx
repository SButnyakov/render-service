import React, { FormEvent, useState } from "react"
import { auth, registration } from "../../http/userAPI"
import { observer } from "mobx-react-lite"

export const SignupForm = observer(() => {
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [email, setEmail] = useState('')

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault()
    signUp()
  }

  const signUp = async () => {
    try {
      const userData = await registration(email, login, password)
      console.log(userData)
    }
    catch (e: any) {
      console.error(e.message)
    }
  }

  return(
    <div >
      <form onSubmit={handleSubmit}>
        <div>
          <label htmlFor="login">Login:</label>
          <input 
            name="login" 
            type="text" 
            value={login} 
            onChange={e => {setLogin(e.target.value)}}
          />
        </div>

        <div>
          <label htmlFor="email">Email:</label>
          <input 
            name="email" 
            type="text" 
            value={email} 
            onChange={e => {setEmail(e.target.value)}}
          />
        </div>

        <div>
          <label htmlFor="password">Password:</label>
          <input 
            name="password" 
            type="password" 
            value={password} 
            onChange={e => {setPassword(e.target.value)}}
          />
        </div>
        <button>Register</button>
      </form>
    </div>
  )
})

export default SignupForm