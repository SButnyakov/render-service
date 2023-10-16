import { Dispatch } from "@reduxjs/toolkit"
import api from "../../api"
import { ILoginRequest} from "../../api/auth/types"
import { loginStart, loginSucess, loginFailure} from "./authReducer"

export const loginUser =
  (data: ILoginRequest) =>
    async (dispatch: Dispatch<any>): Promise<void> => {
      try {
        dispatch(loginStart())

        const res = await api.auth.login(data)

        dispatch(loginSucess(res.data.accessToken))
        
      } catch (e: any) {
        console.error(e)

        dispatch(loginFailure(e.message))
      }
    }

