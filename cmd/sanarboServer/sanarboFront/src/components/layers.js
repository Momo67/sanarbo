export const tile_layers = [
  {
    title: 'Lidar 2016',
    type: 'base',
    layer: 'orthophotos_ortho_lidar_2016',
    url: `https://tiles01.lausanne.ch/tiles/1.0.0/{Layer}/default/2016/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
    requestEncoding: 'REST',
    visible: false,
    zIndex: -1
  },
  {
    title: 'Lidar 2012',
    type: 'base',
    layer: 'orthophotos_ortho_lidar_2012',
    url: `https://tiles01.lausanne.ch/tiles/1.0.0/{Layer}/default/2012/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
    requestEncoding: 'REST',
    visible: false,
    zIndex: -1
  },
  {
    title: 'Plan Ville',
    type: 'base',
    layer: 'plan_ville',
    url: `https://tilesmn95.lausanne.ch/tiles/1.0.0/fonds_geo_osm_bdcad_couleur/default/2021/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
    requestEncoding: 'REST',
    visible: false,
    zIndex: -1
  },
  {
    title: 'Carte nationale',
    type: 'base',
    layer: 'fonds_geo_carte_nationale_msgroup',
    url: `https://tiles01.lausanne.ch/tiles/1.0.0/{Layer}/default/2014/swissgrid_05/{TileMatrix}/{TileRow}/{TileCol}.png`,
    requestEncoding: 'REST',
    visible: false,
    zIndex: -1
  },
];

export const default_tile_grid = {
  origin: [2420000, 1350000],
  resolutions: [50, 20, 10, 5, 2.5, 1, 0.5, 0.25, 0.1, 0.05],
  matrixIds: [0, 1, 2, 3, 4, 5, 6, 7, 8, 9],
};
