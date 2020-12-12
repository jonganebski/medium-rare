import Axios, { AxiosError } from "axios";
import { BASE_URL } from "./constants";
import { delAccountEl } from "./elements.settings";

const deleteAccount = async () => {
  const password = delAccountEl.input?.value;
  if (!password) {
    return;
  }
  const confirmed = confirm(`Your account and every information related with your account will be removed permanently.
Do you want to proceed?`);
  if (!confirmed) {
    return;
  }
  try {
    const { status } = await Axios.delete(BASE_URL + "/api/user", {
      data: { password },
    });
    if (status === 204) {
      document.location.href = "/";
    }
  } catch (err) {
    const errorMessage = (<AxiosError>err).response?.data;
    delAccountEl.desc && (delAccountEl.desc.innerText = errorMessage);
  }
};

const handleDelAccountCancelBtn = (e: Event) => {
  if (!delAccountEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleDelAccountCancelBtn);
  currentTarget?.remove();
  const confirmBtnEl = parentElement?.querySelector(
    ".settings__confirmDelAccount-btn"
  );
  confirmBtnEl?.removeEventListener("click", deleteAccount);
  confirmBtnEl?.remove();
  delAccountEl.desc &&
    (delAccountEl.desc.innerHTML =
      "Your account and all related information will be removed.");
  delAccountEl.desc && (delAccountEl.desc.style.color = "rgba(0, 0, 0, 0.8");
  const delBtn = document.createElement("button");
  delBtn.innerText = "Delete";
  delBtn.className = "settings__gray-btn settings__delAccount-btn";
  delBtn.addEventListener("click", handleDelAccountBtn);
  parentElement?.append(delBtn);
  delAccountEl.input.disabled = true;
  delAccountEl.input.value = "";
  delAccountEl.input.placeholder = "";
};

const handleDelAccountBtn = (e: Event) => {
  if (!delAccountEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleDelAccountBtn);
  currentTarget?.remove();
  const confirmBtnEl = document.createElement("button");
  const cancelBtnEl = document.createElement("button");
  confirmBtnEl.innerText = "DELETE ACCOUNT";
  cancelBtnEl.innerText = "Cancel";
  confirmBtnEl.className = "settings__red-btn settings__confirmDelAccount-btn";
  cancelBtnEl.className = "settings__gray-btn settings__cancelDelAccount-btn";
  confirmBtnEl.addEventListener("click", deleteAccount);
  cancelBtnEl.addEventListener("click", handleDelAccountCancelBtn);
  parentElement?.append(confirmBtnEl);
  parentElement?.append(cancelBtnEl);
  delAccountEl.desc && (delAccountEl.desc.innerHTML = "");
  delAccountEl.desc && (delAccountEl.desc.style.color = "red");
  delAccountEl.input && (delAccountEl.input.disabled = false);
  delAccountEl.input.placeholder = "Enter your password";
};

const initDeleteAccount = () => {
  delAccountEl.delBtn?.addEventListener("click", handleDelAccountBtn);
};

initDeleteAccount();