export const settingsProfile = document.getElementById("settings__profile");

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
