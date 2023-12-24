import React, { useState } from "react";
import cl from '../styles/UploadPageStyles.module.css'
import { sendUploadedFile } from "../http/RenderBackendAPI";

const UploadPage = () => {
  const [isFileUploaded, setIsFileUploaded] = useState(false)
  const [formatSettings, setFormatSettings] = useState('')
  const [resolutionSettings, setResolutionSettings] = useState('')

  const [file, setFile] = useState<string | Blob>('')

  const handleFileUploaded = (event: any) => {
    setIsFileUploaded(true)

    setFile(event.target.files[0])

    console.log(event.target.files)
  }

  const handleSendFile = () => {
    const formData = new FormData()
    formData.append('uploadfile', file)

    try {
      sendUploadedFile(formatSettings, resolutionSettings, formData)
    }
    catch (e) {
      console.error(e)
    }
    
  }

  return(
    <div className="">
      <div className={cl.uploadField}>
        <input type="file" accept=".blend, .glb" onChange={handleFileUploaded}/>
      </div>

      {isFileUploaded && (
        <div>
          <input 
            placeholder="format" 
            type="text" 
            value={formatSettings} 
            onChange={(e) => setFormatSettings(e.target.value)}
          />
          <input 
            placeholder="resolution" 
            type="text" 
            value={resolutionSettings} 
            onChange={(e) => setResolutionSettings(e.target.value)}
          />

          <button onClick={() => handleSendFile()}>
            Send File
          </button>
        </div>
      )}
    </div>
  )
}

export default UploadPage
