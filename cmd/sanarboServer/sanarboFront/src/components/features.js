export const getValidationColor = (idvalidation) => {
  let color = '';
  switch (idvalidation) {
    case 1:   //Existant
      color = '#00FF00';
      break;
    case 5:   //En attente de soins
      color = '#FF00FF';
      break;
    case 6:   //En attente d'abattage
      color = '#FFE600';
      break;
    case 7:   //En attente de remplacement
      color = '#00FFFF';
      break;
    case 8:   //En attente de tomographie
      color = '#0000FF';
      break;
    case 9:   //A surveiller
      color = '#FF0000';
      break;
    case 10:  //En demande d'abattage
      color = '#FF7D00';
      break;
    case 11:  //En attente de projet
      color = '#009696';
      break;
    default:
      color = 'white';
      break;
  };
  return color;
};
