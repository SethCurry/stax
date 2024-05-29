async function get_card_by_name(name) {
  return fetch(
    "https://api.scryfall.com/cards/named?" +
      new URLSearchParams({
        exact: name,
      }),
    {
      method: "GET",
    }
  )
    .then((data) => data.text())
    .then((text) => {
      return JSON.parse(text);
    });
}

async function get_card_image_by_name(name) {
  return get_card_by_name(name)
    .then(({ image_uris }) => {
      return image_uris;
    })
    .then((image_uris) => {
      if (
        image_uris.large !== undefined &&
        image_uris.large !== null &&
        image_uris.large !== ""
      ) {
        return image_uris.large;
      }
      if (
        image_uris.normal !== undefined &&
        image_uris.normal !== null &&
        image_uris.normal !== ""
      ) {
        return image_uris.normal;
      }
      if (
        image_uris.small !== undefined &&
        image_uris.small !== null &&
        image_uris.small !== ""
      ) {
        return image_uris.small;
      }
    });
}
