let isNeedRemoveDuplicates = false

function loadZip(event: Event)  {
  const target = event.target as HTMLInputElement
  if (!target.files || target.files.length != 1) {
    alert("need one file!")
    return
  }
  const file = target.files[0]
  if (file.type != "application/zip") {
    alert("need only .zip!")
  }
}

function main() {

  // TODO

  needRemoveDuplicatesCheckbox.addEventListener("change", () => {
    isNeedRemoveDuplicates = needRemoveDuplicatesCheckbox.checked
  })
  
  // TODO 
  
  const fileInput = document.createElement("input")
  fileInput.type = "file"
  fileInput.style.display = "none"
  fileInput.addEventListener("change", loadZip)
  document.body.appendChild(fileInput)

  const loadButton = document.createElement("button")
  loadButton.textContent = "load .zip"
  loadButton.addEventListener("click", () => {
    fileInput.click()
  })
  document.body.appendChild(loadButton)

}