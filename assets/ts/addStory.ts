const textEditor = document.getElementById("editor__container");
const textEditorTextBoxes = textEditor?.querySelectorAll(".textBox");

const addStory = () => {
  textEditor?.addEventListener("keydown", (e) => {
    const { key } = <KeyboardEvent>e;
    const { innerHTML } = <HTMLParagraphElement>e.currentTarget;
    console.log(key);
    const x = innerHTML + key;

    textEditor.innerHTML = x;
  });
};

addStory();
