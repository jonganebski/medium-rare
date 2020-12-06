import { useEditor } from "./useEditor";

const addStory = () => {
  if (document.location.pathname.includes("new-story")) {
    const blocks = [
      {
        type: "header",
        data: { level: 2, text: "Title" },
      },
      {
        type: "paragraph",
        data: { text: "Write your story" },
      },
    ];
    useEditor("editor-write", "Write your story", blocks);
  }
};

addStory();
