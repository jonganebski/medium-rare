import Axios from "axios";
import { BASE_URL } from "./constants";
import { editBioEl } from "./elements.settings";

let prevBioInputValue = "";

const updateBio = async (e: Event) => {
  if (!editBioEl.input) {
    return;
  }
  const newBio = editBioEl.input?.value;
  const { status, data: updatedBio } = await Axios.patch(
    BASE_URL + `/api/user/bio`,
    {
      bio: newBio,
    }
  );
  if (status === 200) {
    const target = e.target as HTMLButtonElement | null;
    const parentElement = target?.parentElement as HTMLElement | null;
    target?.removeEventListener("click", updateBio);
    target?.remove();
    const cancelBtnEl = parentElement?.querySelector(
      ".settings__cancelBio-btn"
    );
    cancelBtnEl?.removeEventListener("click", handleBioEditCancelBtn);
    cancelBtnEl?.remove();
    const editBtnEl = document.createElement("button");
    editBtnEl.innerText = "Edit";
    editBtnEl.className = "settings__gray-btn settings__editBio-btn";
    editBtnEl.addEventListener("click", handleBioEditBtn);
    parentElement?.append(editBtnEl);
    editBioEl.input.disabled = true;
    prevBioInputValue = updatedBio;
  }
};

const handleBioEditCancelBtn = (e: Event) => {
  if (!editBioEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleBioEditCancelBtn);
  currentTarget?.remove();
  const saveBtnEl = parentElement?.querySelector(".settings__saveBio-btn");
  saveBtnEl?.removeEventListener("click", updateBio);
  saveBtnEl?.remove();
  const editBtnEl = document.createElement("button");
  editBtnEl.innerText = "Edit";
  editBtnEl.className = "settings__gray-btn settings__editBio-btn";
  editBtnEl.addEventListener("click", handleBioEditBtn);
  parentElement?.append(editBtnEl);
  editBioEl.input.disabled = true;
  editBioEl.input.value = prevBioInputValue;
};

const handleBioEditBtn = (e: Event) => {
  if (!editBioEl.input) {
    return;
  }
  const currentTarget = e.currentTarget as HTMLButtonElement | null;
  const parentElement = currentTarget?.parentElement as HTMLElement | null;
  currentTarget?.removeEventListener("click", handleBioEditBtn);
  currentTarget?.remove();
  const saveBtnEl = document.createElement("button");
  const cancelBtnEl = document.createElement("button");
  saveBtnEl.innerText = "Save";
  cancelBtnEl.innerText = "Cancel";
  saveBtnEl.className = "settings__green-btn settings__saveBio-btn";
  cancelBtnEl.className = "settings__gray-btn settings__cancelBio-btn";
  saveBtnEl.addEventListener("click", updateBio);
  cancelBtnEl.addEventListener("click", handleBioEditCancelBtn);
  parentElement?.append(saveBtnEl);
  parentElement?.append(cancelBtnEl);
  prevBioInputValue = editBioEl.input.value;
  editBioEl.input && (editBioEl.input.disabled = false);
};

const useSettings = () => {
  editBioEl.editBtn?.addEventListener("click", handleBioEditBtn);
};

useSettings();
