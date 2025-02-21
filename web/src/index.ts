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
  const needRemoveDuplicatesCheckbox = document.createElement("input")
  needRemoveDuplicatesCheckbox.type = "checkbox"
  needRemoveDuplicatesCheckbox.id = "toggleSwitch"

  const needRemoveDuplicatesLabel = document.createElement("label")
  needRemoveDuplicatesLabel.appendChild(needRemoveDuplicatesCheckbox)
  needRemoveDuplicatesLabel.appendChild(document.createTextNode("remove duplicates?"))

  needRemoveDuplicatesCheckbox.addEventListener("change", () => {
    isNeedRemoveDuplicates = needRemoveDuplicatesCheckbox.checked
  })

  document.body.appendChild(needRemoveDuplicatesLabel)

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

  document.body.appendChild(document.createElement("hr"))

  const author = document.createElement("a")
  author.href = "https://github.com/Stasenko-Konstantin"
  author.textContent = "author github page"
  author.target = "_blank"
  document.body.appendChild(author)
}

main()