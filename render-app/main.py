import subprocess
import time
import os

import requests

request_url = "http://localhost:8082/request"

FIVE_MINUTES = 300

SCRIPT_PATH = "render.py"
BLENDER_PATH = "E:/blenda/blender.exe"


def run_blender(blender: str, project: str, script: str) -> None:
    output = subprocess.check_output([blender,
                                      project,
                                      '--background',
                                      '--python', script,
                                      '--render-output', '//frame',
                                      '--render-frame', '0'])
    print(output.decode('utf-8'))


def download_file(download_link: str, filename: str) -> bool:
    response = requests.get(download_link)

    if response.status_code != 200:
        return False

    open(filename, 'wb').write(response.content)

    return True


def parse_download_link(download_link: str) -> (str, str):
    arr = download_link.split("/")
    return arr[len(arr)-4], arr[len(arr)-1]


def main() -> None:
    while True:
        response = requests.get(request_url)
        print("requesting new order")

        if response.status_code != 200:
            print("error request")
            time.sleep(FIVE_MINUTES)
            continue

        json = response.json()

        if json["status"] == "Empty":
            print("empty")
            time.sleep(20)
            continue

        print("got new order")

        print("downloading file")
        uid, filename = parse_download_link(json["download_link"])
        link_filename = filename.split(".")[0]

        downloaded = download_file(json["download_link"], filename)
        if not downloaded:
            print("error downloading")
            requests.put(f"http://localhost:8082/{uid}/blend/update/{link_filename}/error")
            continue

        print("file downloaded")

        print("updating status: IN PROGRESS")
        requests.put(f"http://localhost:8082/{uid}/blend/update/{link_filename}/in-progress")

        print("running blender")
        run_blender(BLENDER_PATH, filename, SCRIPT_PATH)
        print("render finished")

        image_name = f"{link_filename}.png"

        os.rename("frame0000.png", image_name)

        files = {'uploadfile': open(image_name, 'rb')}

        requests.post(url=f"http://localhost:8082/{uid}/image/upload", files=files)

        #os.remove(filename)
        #os.remove(image_name)


if __name__ == '__main__':
    main()
