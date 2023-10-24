import React from "react";

export const UploadPage = () => {
  return(
    <div>
      <form method="post">
        <input id="file" name="file" type="file" />
        <button>Upload</button>
      </form>
    </div>
  )
}

export default UploadPage;