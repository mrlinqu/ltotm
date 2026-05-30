function showToast(msg) {
  const toast = document.getElementById("toast");
  toast.textContent = msg;
  toast.classList.add("show");
  setTimeout(() => toast.classList.remove("show"), 2500);
}

function showError(msg) {
  const toast_err = document.getElementById("toastErr");
  toast_err.textContent = msg;
  toast_err.classList.add("show");
  toast_err.onclick = () => toast_err.classList.remove("show");
}

let selectedHours = 24;

const textInput = document.getElementById("textInput");
const charCountEl = document.getElementById("charCount");

textInput.addEventListener("input", () => {
  charCountEl.textContent = textInput.value.length.toLocaleString("ru-RU");
});

function selectDuration(btn) {
  document
    .querySelectorAll(".duration-btn")
    .forEach((b) => b.classList.remove("active"));
  btn.classList.add("active");
  selectedHours = parseInt(btn.dataset.hours);
}

function generateKey() {
  var res = [];
  for (var i = 16; i > 0; i--) res.push(Math.round((Math.random() * 1000) / 4));

  return res;
}

function encrypt(msg, keyBytes) {
  const textBytes = aesjs.utils.utf8.toBytes(msg);
  const aesCtr = new aesjs.ModeOfOperation.ctr(keyBytes, new aesjs.Counter(5));
  const encryptedBytes = aesCtr.encrypt(textBytes);

  return encryptedBytes;
}

function decrypt(encryptedBytes, keyBytes) {
  const aesCtr = new aesjs.ModeOfOperation.ctr(keyBytes, new aesjs.Counter(5));
  const msgBytes = aesCtr.decrypt(encryptedBytes);
  const msg = aesjs.utils.utf8.fromBytes(msgBytes);

  return msg;
}

async function put(msg, ttl) {
  const resp = await fetch("/put", {
    method: "POST",
    body: JSON.stringify({ msg, ttl }),
  });

  const json = await resp.json();

  return json;
}

async function get(id) {
  const resp = await fetch("/get", {
    method: "POST",
    body: JSON.stringify({ id }),
  });

  const json = await resp.json();

  return json;
}

function generateId() {
  const chars =
    "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789";
  let id = "";
  for (let i = 0; i < 10; i++) {
    id += chars.charAt(Math.floor(Math.random() * chars.length));
  }
  return id;
}

function disableSendBtn() {
  const btn = document.getElementById("submitBtn");
  btn.disabled = true;
  btn.innerHTML = '<span class="spinner"></span>Сохраняем…';
}

function enableSendBtn() {
  const btn = document.getElementById("submitBtn");
  btn.disabled = false;
  btn.innerHTML = "Создать ссылку";
}

function disableReceiveBtn() {
  const btn = document.getElementById("receiveBtn");
  btn.disabled = true;
  btn.innerHTML = '<span class="spinner"></span>Получаем…';
}

function enableReceiveBtn() {
  const btn = document.getElementById("receiveBtn");
  btn.disabled = false;
  btn.innerHTML = "Получить сообщение";
}

function createLink() {
  const msg = textInput.value.trim();
  if (!msg) {
    showToast("Введите сообщение");
    textInput.focus();
    return;
  }

  disableSendBtn();

  const key = generateKey();
  const cryptedMsg = encrypt(msg, key);
  const keyUrl = fromByteArray(key);

  put(fromByteArray(cryptedMsg), selectedHours)
    .then((data) => {
      enableSendBtn();
      let generatedUrl = window.location.origin + "#" + data.id + keyUrl;
      showLink(generatedUrl, new Date(data.expiration_time));
    })
    .catch((e) => {
      console.log(e);
      showError("Ошибка! Не удалось сохранить сообщение!");
      enableSendBtn();
    });
}

function showLink(link, expirationTime) {
  document.getElementById("resultLink").textContent = link;
  document.getElementById("expiryDate").textContent =
    expirationTime.toLocaleDateString("ru-RU", {
      day: "numeric",
      month: "long",
      year: "numeric",
      hour: "2-digit",
      minute: "2-digit",
    });

  activateForm("resultView");
  showToast("Ссылка создана!");
}

function copyLink() {
  const copyBtn = document.getElementById("copyBtn");
  const copyIcon = document.getElementById("copyIcon");
  const checkIcon = document.getElementById("checkIcon");
  const copyBtnText = document.getElementById("copyBtnText");
  const generatedUrl = document.getElementById("resultLink").textContent;

  navigator.clipboard.writeText(generatedUrl).then(() => {
    copyBtn.classList.add("copied");
    copyIcon.style.display = "none";
    checkIcon.style.display = "block";
    copyBtnText.textContent = "Скопировано";

    setTimeout(() => {
      copyBtn.classList.remove("copied");
      copyIcon.style.display = "block";
      checkIcon.style.display = "none";
      copyBtnText.textContent = "Копировать";
    }, 2000);
  });
}

function copyMsg() {
  const copyBtn = document.getElementById("copyMsgBtn");
  const copyIcon = document.getElementById("copyMsgIcon");
  const checkIcon = document.getElementById("checkMsgIcon");
  const copyBtnText = document.getElementById("copyMsgBtnText");
  const result = document.getElementById("resultMessage").textContent;

  navigator.clipboard.writeText(result).then(() => {
    copyBtn.classList.add("copied");
    copyIcon.style.display = "none";
    checkIcon.style.display = "block";
    copyBtnText.textContent = "Скопировано";

    setTimeout(() => {
      copyBtn.classList.remove("copied");
      copyIcon.style.display = "block";
      checkIcon.style.display = "none";
      copyBtnText.textContent = "Копировать";
    }, 2000);
  });
}

function receiveMessage() {
  const id = window.location.hash.substring(1, 25);
  const key = window.location.hash.substring(25);

  disableReceiveBtn();

  get(id)
    .then((data) => {
      const msg = decrypt(toByteArray(data.msg), toByteArray(key));

      document.getElementById("resultMessage").textContent = msg;

      activateForm("receivedMessage");
      enableReceiveBtn();
    })
    .catch((e) => {
      console.log(e);
      showError("Ошибка! Не удалось получить сообщение!");
      enableReceiveBtn();
    });
}

function resetForm() {
  textInput.value = "";
  charCountEl.textContent = "0";

  activateForm("formView");
  document.getElementById("formView").style.animation = "none";
  void document.getElementById("formView").offsetHeight;
  document.getElementById("formView").style.animation =
    "fadeSlideIn 0.4s ease-out";

  textInput.focus();
}

function hideForm(formId) {
  document.getElementById(formId).classList.remove("visible");
  document.getElementById(formId).classList.add("hidden");
}

function showForm(formId) {
  document.getElementById(formId).classList.remove("hidden");
  document.getElementById(formId).classList.add("visible");
}

function activateForm(formId) {
  hideForm("formView");
  hideForm("resultView");
  hideForm("receiveMessage");
  hideForm("receivedMessage");

  showForm(formId);
}

if (window.location.hash != "") {
  showForm("receiveMessage");
} else {
  resetForm();
}
