import React, { FormEvent, useState } from "react";
import { useAppDispatch } from "../../store";
import { loginUser } from "../../store/auth/actionCreators";

export const SignInForm = () => {
  const [login, setLogin] = useState('');
  const [password, setPassword] = useState('');

  const disapatch = useAppDispatch()

  const handleSubmit = (e: FormEvent) => {
    e.preventDefault();
    disapatch(loginUser({login, password}));
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
            type="text" 
            value={password} 
            onChange={e => {setPassword(e.target.value)}}
          />
        </div>
        <button>LogIn</button>
      </form>
    </div>
  )
}

export default SignInForm;