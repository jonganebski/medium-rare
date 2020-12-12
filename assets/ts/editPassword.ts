import Axios from "axios";
import { BASE_URL } from "./constants";
import { editPasswordEl, settingsSecurity } from "./elements.settings";

const handlePassEditCancelBtn = (e: Event) => {
  if (!editPasswordEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handlePassEditCancelBtn);
  currentTarget?.remove();
  const saveBtnEl = parentElement?.querySelector(".settings__savePass-btn");
  saveBtnEl?.removeEventListener("click", updatePassword);
  saveBtnEl?.remove();
  const firstInputEl = settingsSecurity?.querySelector(
    ".settings__passInput-1"
  );
  const secondInputEl = settingsSecurity?.querySelector(
    ".settings__passInput-2"
  );
  firstInputEl?.remove();
  secondInputEl?.remove();
  editPasswordEl.desc &&
    (editPasswordEl.desc.innerHTML = "Change your password.");
  editPasswordEl.desc &&
    (editPasswordEl.desc.style.color = "rgba(0, 0, 0, 0.8");
  const editBtnEl = document.createElement("button");
  editBtnEl.innerText = "Edit";
  editBtnEl.className = "settings__gray-btn settings__editPass-btn";
  editBtnEl.addEventListener("click", handlePassEditBtn);
  parentElement?.append(editBtnEl);
  editPasswordEl.input.disabled = true;
  editPasswordEl.input.value = "";
  editPasswordEl.input.placeholder = "";
};

const handlePassEditBtn = (e: Event) => {
  if (!editPasswordEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", () => {});
  currentTarget?.remove();
  const saveBtnEl = document.createElement("button");
  const cancelBtnEl = document.createElement("button");
  const firstInputEl = document.createElement("input");
  const secondInputEl = document.createElement("input");
  saveBtnEl.innerText = "Save";
  cancelBtnEl.innerText = "Cancel";
  saveBtnEl.className = "settings__green-btn settings__savePass-btn";
  cancelBtnEl.className = "settings__gray-btn settings__cancelPass-btn";
  firstInputEl.className = "settings__stack-input settings__passInput-1";
  secondInputEl.className = "settings__stack-input settings__passInput-2";
  firstInputEl.type = "password";
  secondInputEl.type = "password";
  firstInputEl.placeholder = "Enter new password";
  secondInputEl.placeholder = "Enter new password again";
  saveBtnEl.addEventListener("click", updatePassword);
  cancelBtnEl.addEventListener("click", handlePassEditCancelBtn);
  parentElement?.append(saveBtnEl);
  parentElement?.append(cancelBtnEl);
  editPasswordEl.input.parentElement?.append(firstInputEl);
  editPasswordEl.input.parentElement?.append(secondInputEl);
  editPasswordEl.desc && (editPasswordEl.desc.innerHTML = "");
  editPasswordEl.desc && (editPasswordEl.desc.style.color = "red");
  editPasswordEl.input && (editPasswordEl.input.disabled = false);
  editPasswordEl.input.placeholder = "Enter your current password";
};

const updatePassword = async (e: Event) => {
  const firstInputEl = settingsSecurity?.querySelector(
    ".settings__passInput-1"
  ) as HTMLInputElement | null;
  const secondInputEl = settingsSecurity?.querySelector(
    ".settings__passInput-2"
  ) as HTMLInputElement | null;
  const originalPass = editPasswordEl.input?.value;
  const firstPass = firstInputEl?.value;
  const secondPass = secondInputEl?.value;
  if (!originalPass || !firstPass || !secondPass) {
    editPasswordEl.desc &&
      (editPasswordEl.desc.innerHTML = "Please enter passwords");
    return;
  }
  if (firstPass.length < 6 || secondPass.length < 6) {
    editPasswordEl.desc &&
      (editPasswordEl.desc.innerHTML =
        "Password must longer than 5 characters");
    return;
  }
  if (firstPass !== secondPass) {
    editPasswordEl.desc &&
      (editPasswordEl.desc.innerHTML = "Passwords don't match");
    return;
  }
  const isConfirmed = confirm(
    "You are changing password of your accound. Proceed?"
  );
  if (!isConfirmed) {
    return;
  }
  try {
    const { status } = await Axios.patch(BASE_URL + "/api/user/password", {
      originalPass,
      firstPass,
      secondPass,
    });
    if (status === 200) {
      if (!editPasswordEl.input) {
        return;
      }
      const currentTarget = e.currentTarget as HTMLButtonElement | null;
      const parentElement = currentTarget?.parentElement as HTMLElement | null;
      currentTarget?.removeEventListener("click", updatePassword);
      currentTarget?.remove();
      const cancelBtnEl = parentElement?.querySelector(
        ".settings__cancelPass-btn"
      );
      cancelBtnEl?.removeEventListener("click", handlePassEditCancelBtn);
      cancelBtnEl?.remove();
      firstInputEl?.remove();
      secondInputEl?.remove();
      editPasswordEl.desc &&
        (editPasswordEl.desc.innerHTML = "Change your password.");
      editPasswordEl.desc &&
        (editPasswordEl.desc.style.color = "rgba(0, 0, 0, 0.8");
      const editBtnEl = document.createElement("button");
      editBtnEl.innerText = "Edit";
      editBtnEl.className = "settings__gray-btn settings__editPass-btn";
      editBtnEl.addEventListener("click", handlePassEditBtn);
      parentElement?.append(editBtnEl);
      editPasswordEl.input.disabled = true;
      editPasswordEl.input.value = "";
      editPasswordEl.input.placeholder = "Password successfully changed.";
    }
  } catch {
    editPasswordEl.desc &&
      (editPasswordEl.desc.innerHTML = "Something's wrong.");
  }
};

const initEditPassword = () => {
  editPasswordEl.editBtn?.addEventListener("click", handlePassEditBtn);
};

initEditPassword();
