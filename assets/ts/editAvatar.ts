import Axios from "axios";
import { BASE_URL } from "./constants";
import { avatarFrame } from "./elements.header";
import { editAvatarEl, settingsProfile } from "./elements.settings";

let prevAvatarInputValue = "";

const handleAvatarEditCancelBtn = (e: Event) => {
  if (!editAvatarEl.input || !editAvatarEl.avatar) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleAvatarEditCancelBtn);
  currentTarget?.remove();
  const saveBtnEl = parentElement?.querySelector(".settings__saveAvatar-btn");
  saveBtnEl?.remove();
  const guideEl = settingsProfile?.querySelector(
    ".settings__uploadAvatar-guide"
  );
  guideEl && (guideEl.className = "settings__uploadAvatar-guide _flex-cc");
  guideEl?.removeEventListener("click", handleAvatarClick);
  guideEl?.remove();
  const editBtnEl = document.createElement("button");
  editBtnEl.innerText = "Edit";
  editBtnEl.className = "settings__gray-btn settings__editAvatar-btn";
  editBtnEl.addEventListener("click", handleAvatarEditBtn);
  editAvatarEl.avatar.removeEventListener("click", handleAvatarClick);
  parentElement?.append(editBtnEl);
  editAvatarEl.input.disabled = true;
  editAvatarEl.avatar.src = prevAvatarInputValue;
};

const handleAvatarClick = () => {
  editAvatarEl.input?.click();
};

const updateAvatar = async (e: Event) => {
  e.preventDefault();
  const form = e.currentTarget as HTMLFormElement | null;
  if (!form) {
    return;
  }
  if (editAvatarEl.avatar?.src == prevAvatarInputValue) {
    return;
  }
  const formData = new FormData(form);
  formData.append("oldAvatarUrl", prevAvatarInputValue);
  const { status, data: newUrl } = await Axios.patch(
    BASE_URL + "/api/user/avatar",
    formData
  );
  if (status === 200) {
    const guideEl = settingsProfile?.querySelector(
      ".settings__uploadAvatar-guide"
    );
    guideEl && (guideEl.className = "settings__uploadAvatar-guide _flex-cc");
    guideEl?.removeEventListener("click", handleAvatarClick);
    guideEl?.remove();
    const cancelBtnEl = form?.querySelector(".settings__cancelAvatar-btn");
    cancelBtnEl?.removeEventListener("click", handleAvatarEditCancelBtn);
    cancelBtnEl?.remove();
    const saveBtn = form?.querySelector(".settings__saveAvatar-btn");
    saveBtn?.remove();
    const editBtnEl = document.createElement("button");
    editBtnEl.innerText = "Edit";
    editBtnEl.className = "settings__gray-btn settings__editAvatar-btn";
    editBtnEl.addEventListener("click", handleAvatarEditBtn);
    const controlBox = form.querySelector(".settings__stack-control");
    controlBox?.append(editBtnEl);
    editAvatarEl.input!.disabled = true;
    editAvatarEl.input!.value = "";

    avatarFrame!.querySelector("img")!.src = newUrl;
    editAvatarEl.avatar!.src = newUrl;
    prevAvatarInputValue = newUrl;
  }
};

const handleAvatarEditBtn = (e: Event) => {
  if (!editAvatarEl.input || !editAvatarEl.avatar) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleAvatarEditBtn);
  currentTarget?.remove();
  const saveBtnEl = document.createElement("button");
  const cancelBtnEl = document.createElement("button");
  saveBtnEl.type = "submit";
  saveBtnEl.innerText = "Save";
  cancelBtnEl.innerText = "Cancel";
  saveBtnEl.className = "settings__green-btn settings__saveAvatar-btn";
  cancelBtnEl.className = "settings__gray-btn settings__cancelAvatar-btn";
  cancelBtnEl.addEventListener("click", handleAvatarEditCancelBtn);
  const guideEl = document.createElement("div");
  guideEl.innerText = "Click here";
  guideEl.className = "settings__uploadAvatar-guide _flex-cc";
  guideEl.addEventListener("click", handleAvatarClick);
  editAvatarEl.avatar.parentElement?.append(guideEl);
  parentElement?.append(saveBtnEl);
  parentElement?.append(cancelBtnEl);
  prevAvatarInputValue = editAvatarEl.avatar.src;
  editAvatarEl.input && (editAvatarEl.input.disabled = false);
};

const onFileInput = (e: Event) => {
  const files = (<HTMLInputElement>e.currentTarget).files;
  if (!files) {
    return;
  }
  const reader = new FileReader();
  const file = files[0];
  reader.addEventListener("loadend", () => {
    if (!reader.result?.toString() || !editAvatarEl.avatar) {
      return;
    }
    editAvatarEl.avatar.src = reader.result?.toString();
    const guideEl = settingsProfile?.querySelector(
      ".settings__uploadAvatar-guide"
    ) as HTMLDivElement | null;
    guideEl && (guideEl.style.backgroundColor = "transparent");
    guideEl && (guideEl.innerText = "");
  });
  reader.readAsDataURL(file);
};

const initEditAvatar = () => {
  editAvatarEl.editBtn?.addEventListener("click", handleAvatarEditBtn);
  editAvatarEl.input?.addEventListener("input", onFileInput);
  editAvatarEl.form?.addEventListener("submit", updateAvatar);
};

initEditAvatar();
