export const settingsProfile = document.getElementById("settings__profile");
export const settingsSecurity = document.getElementById("settings__security");

export const editUsernameEl = {
  input: settingsProfile?.querySelector(
    ".settings__usernameInput"
  ) as HTMLInputElement | null,
  editBtn: settingsProfile?.querySelector(".settings__editUsername-btn"),
};

export const editBioEl = {
  input: settingsProfile?.querySelector(
    ".settings__bioInput"
  ) as HTMLInputElement | null,
  editBtn: settingsProfile?.querySelector(".settings__editBio-btn"),
};

export const editAvatarEl = {
  form: settingsProfile?.querySelector(".settings__stack-avatar-form"),
  input: settingsProfile?.querySelector(
    ".settings__avatarInput"
  ) as HTMLInputElement | null,
  avatar: settingsProfile?.querySelector(
    ".settings__avatar-img"
  ) as HTMLImageElement | null,
  editBtn: settingsProfile?.querySelector(".settings__editAvatar-btn"),
};

export const editPasswordEl = {
  input: settingsSecurity?.querySelector(
    ".settings__passInput"
  ) as HTMLInputElement | null,
  editBtn: settingsSecurity?.querySelector(".settings__editPass-btn"),
  desc: settingsSecurity?.querySelector(
    ".settings__passDesc"
  ) as HTMLElement | null,
};
