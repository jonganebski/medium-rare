import Axios from "axios";
import { BASE_URL } from "./constants";
import { editUsernameEl } from "./elements.settings";

let prevUsernameInputValue = "";

const updateUsername = async (e: Event) => {
  console.log("foo");
  if (!editUsernameEl.input) {
    return;
  }
  const newUsername = editUsernameEl.input?.value;
  if (!newUsername) {
    return;
  }
  const { status, data: updatedUsername } = await Axios.patch(
    BASE_URL + `/api/user/username`,
    {
      username: newUsername,
    }
  );
  if (status === 200) {
    const target = e.target as HTMLButtonElement | null;
    const parentElement = target?.parentElement as HTMLElement | null;
    target?.removeEventListener("click", updateUsername);
    target?.remove();
    const cancelBtnEl = parentElement?.querySelector(
      ".settings__cancelUsername-btn"
    );
    cancelBtnEl?.removeEventListener("click", handleUsernameEditCancelBtn);
    cancelBtnEl?.remove();
    const editBtnEl = document.createElement("button");
    editBtnEl.innerText = "Edit";
    editBtnEl.className = "settings__gray-btn settings__editUsername-btn";
    editBtnEl.addEventListener("click", handleUsernameEditBtn);
    parentElement?.append(editBtnEl);
    editUsernameEl.input.disabled = true;
    prevUsernameInputValue = updatedUsername;
  }
};

const handleUsernameEditCancelBtn = (e: Event) => {
  if (!editUsernameEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleUsernameEditCancelBtn);
  currentTarget?.remove();
  const saveBtnEl = parentElement?.querySelector(".settings__saveUsername-btn");
  saveBtnEl?.removeEventListener("click", updateUsername);
  saveBtnEl?.remove();
  const editBtnEl = document.createElement("button");
  editBtnEl.innerText = "Edit";
  editBtnEl.className = "settings__gray-btn settings__editUsername-btn";
  editBtnEl.addEventListener("click", handleUsernameEditBtn);
  parentElement?.append(editBtnEl);
  editUsernameEl.input.disabled = true;
  editUsernameEl.input.value = prevUsernameInputValue;
};

const handleUsernameEditBtn = (e: Event) => {
  if (!editUsernameEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleUsernameEditBtn);
  currentTarget?.remove();
  const saveBtnEl = document.createElement("button");
  const cancelBtnEl = document.createElement("button");
  saveBtnEl.innerText = "Save";
  cancelBtnEl.innerText = "Cancel";
  saveBtnEl.className = "settings__green-btn settings__saveUsername-btn";
  cancelBtnEl.className = "settings__gray-btn settings__cancelUsername-btn";
  saveBtnEl.addEventListener("click", updateUsername);
  cancelBtnEl.addEventListener("click", handleUsernameEditCancelBtn);
  parentElement?.append(saveBtnEl);
  parentElement?.append(cancelBtnEl);
  prevUsernameInputValue = editUsernameEl.input.value;
  editUsernameEl.input && (editUsernameEl.input.disabled = false);
};

const useSettings = () => {
  editUsernameEl.editBtn?.addEventListener("click", handleUsernameEditBtn);
};

useSettings();
