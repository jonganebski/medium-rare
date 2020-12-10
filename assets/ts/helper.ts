export const getIdParam = (divider: string) => {
  const splitedPath = document.location.pathname.split(divider);
  const idParam = splitedPath[1].replace(/[/]/g, "");
  return idParam;
};
