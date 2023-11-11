import React, { FormEvent, useState } from "react";
import { auth } from "../../http/userAPI";
import { observer } from "mobx-react-lite";

import UserStore from "../../store/UserStore";

export const SignInForm = observer(() => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    signIn()
  }

  const signIn = async () => {
    try {
      const userData = await auth(login, password)

      UserStore.setUser(userData)
      UserStore.setIsAuth(true)
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
          <label htmlFor="password">Password:</label>
          <input 
            name="password" 
            type="password" 
            value={password} 
            onChange={e => {setPassword(e.target.value)}}
          />
        </div>
        <button>LogIn</button>
      </form>
    </div>
  )
})

export default SignInForm;