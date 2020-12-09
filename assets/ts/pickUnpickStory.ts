import Axios from "axios";
import { BASE_URL } from "./constants";
import { pickStoryBtn, unpickStoryBtn } from "./elements.readStory";

const pickStory = async (e: Event) => {
  const splitedPath = document.location.pathname.split("read");
  const storyId = splitedPath[1].replace(/[/]/g, "");
  const { status } = await Axios.post(BASE_URL + `/admin/pick/${storyId}`);
  if (status === 200) {
    const target = e.target as HTMLButtonElement | null;
    if (target) {
      target.className = "editor-readOnly__unpick";
      target.innerText = "Unpick";
      target.removeEventListener("click", pickStory);
      target.addEventListener("click", unpickStory);
    }
  }
};

const unpickStory = async (e: Event) => {
  const splitedPath = document.location.pathname.split("read");
  const storyId = splitedPath[1].replace(/[/]/g, "");
  const { status } = await Axios.post(BASE_URL + `/admin/unpick/${storyId}`);
  if (status === 200) {
    const target = e.target as HTMLButtonElement | null;
    if (target) {
      target.className = "editor-readOnly__pick";
      target.innerText = "Pick";
      target.removeEventListener("click", unpickStory);
      target.addEventListener("click", pickStory);
    }
  }
};

const initPickUnpickStory = () => {
  pickStoryBtn?.addEventListener("click", pickStory);
  unpickStoryBtn?.addEventListener("click", unpickStory);
};

initPickUnpickStory();
