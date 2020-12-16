import Axios from "axios";
import { pickStoryBtn, unpickStoryBtn } from "./elements.readStory";
import { getIdParam } from "./helper";

const pickStory = async (e: Event) => {
  const storyId = getIdParam("read-story");
  try {
    const { status } = await Axios.patch(`/api/admin/pick/${storyId}`);
    if (status < 300) {
      const target = e.target as HTMLButtonElement | null;
      if (target) {
        target.className = "editor-readOnly__unpick";
        target.innerText = "Unpick";
        target.removeEventListener("click", pickStory);
        target.addEventListener("click", unpickStory);
      }
    }
  } catch {
    alert("Failed to pick. Please try again.");
  }
};

const unpickStory = async (e: Event) => {
  const storyId = getIdParam("read-story");
  try {
    const { status } = await Axios.post(`/api/admin/unpick/${storyId}`);
    if (status < 300) {
      const target = e.target as HTMLButtonElement | null;
      if (target) {
        target.className = "editor-readOnly__pick";
        target.innerText = "Pick";
        target.removeEventListener("click", unpickStory);
        target.addEventListener("click", pickStory);
      }
    }
  } catch {
    alert("Failed to unpick. Please try again.");
  }
};

const init = () => {
  pickStoryBtn?.addEventListener("click", pickStory);
  unpickStoryBtn?.addEventListener("click", unpickStory);
};

init();
