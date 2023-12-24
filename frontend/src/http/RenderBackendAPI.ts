import { $authHost } from "."

/* TODO: Пока в body ручки можно отправлять только сам файл FormData. Позже нужно будет переделать */ 
export const sendUploadedFile = async (format: string, resolution: string, file: FormData) => {
  try {
    const {data} = await $authHost.post('senfd', file, {headers: {
      'Content-Type': 'multipart/form-data',
    }})
  }
  catch (e) {
    console.error(e)
  }
}
