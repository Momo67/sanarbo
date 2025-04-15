import 'ol/ol.css'
import OlMap from 'ol/Map';
import OlView from 'ol/View';
import OlLayerTile from 'ol/layer/Tile';
import OlSourceWMTS, {optionsFromCapabilities} from 'ol/source/WMTS';
import OlFormatWMTSCapabilities from 'ol/format/WMTSCapabilities';
import OlProjection from 'ol/proj/Projection';
import proj4 from 'proj4';
import {register} from 'ol/proj/proj4';


const lausanneGareinMN95 = [2537968.5, 1152088.0];
const defaultBaseLayer = 'fonds_geo_osm_bdcad_couleur';

const urlSWISSTOPO = 'https://wmts.geo.admin.ch/EPSG/2056/1.0.0/WMTSCapabilities.xml?lang=fr';
const urlLausanneMN95 = 'https://tilesmn95.lausanne.ch/tiles/1.0.0/LausanneWMTS.xml';

proj4.defs(
  'EPSG:2056',
  '+proj=somerc +lat_0=46.95240555555556 +lon_0=7.439583333333333 +k_0=1 +x_0=2600000 +y_0=1200000 +ellps=bessel +towgs84=674.374,15.056,405.346,0,0,0,0 +units=m +no_defs'
);
//const LausanneLonLat = [6.62925, 46.51735];
const MaxExtent = [2532500, 1149000, 2545625, 1161000];
register(proj4);
const swissProjection = new OlProjection({
		code: 'EPSG:2056',
		extent: MaxExtent,
		units: 'm',
});
const parser = new OlFormatWMTSCapabilities();
//const isNullOrUndefined = (v) => typeof variable === 'undefined' || variable === null || variable === '';
const setBaseLayer = (olMap, baseLayerName) => {
		console.log(`## in setBaseLayer(${baseLayerName})`);
		const localDebug = false;
		let isBaseLayerNameFound = false;
		olMap.getLayers()
		     .forEach((layer) => {
				     const type = layer.get('type');
				     const source = layer.getSource();
				     if (type === 'base') {
						     const currentBaseLayer = source.getLayer();
						     console.log(`currentBaseLayer : ${currentBaseLayer}`)
						     if (currentBaseLayer === baseLayerName) {
								     layer.setVisible(true);
								     isBaseLayerNameFound = true;
						     } else {
								     layer.setVisible(false);
						     }
						     if (localDebug) {
								     console.log(`layers : ${currentBaseLayer} [${type}]
                 title : ${layer.get('title')}
                 isvisible: ${layer.get('visible')}`, layer, source);
						     }
				     }
		      });
		if (!isBaseLayerNameFound) {
				console.log(`WARNING : The layer ${baseLayerName} was not found !`);
		}
};


function getWmtsSource (WMTSCapabilitiesParsed, layerName) {
		const localDebug = false;
		if (localDebug) console.log(`layerName: ${layerName}`);
		const WMTSOptions = optionsFromCapabilities(WMTSCapabilitiesParsed, {
				layer: layerName,
				matrixSet: 'EPSG2056',
				format: 'image/png',
				style: 'default',
				crossOrigin: 'anonymous',
		});
		return new OlSourceWMTS(WMTSOptions);
}

function createBaseOlLayerTile (parsedWmtsCapabilities, title, layerName, visible = false) {
		return new OlLayerTile({
				title,
				type: 'base',
				visible,
				source: getWmtsSource(parsedWmtsCapabilities, layerName),
		});
}

async function getWMTSCapabilitiesFromUrl (url) {
		const response = await fetch(url);
		if (!response.ok) {
				const message = `###!### ERROR in getWMTSCapabilitiesFromUrl when doing fetch(${url}: http status: ${response.status}`;
				throw new Error(message);
		}
		const WMTSCapabilities = await response.text();
		return WMTSCapabilities;
}

async function getWmtsBaseLayers (url) {
  const arrWmtsLayers = [];
  try {
      const WMTSCapabilities = await getWMTSCapabilitiesFromUrl(url);
      const WMTSCapabilitiesParsed = parser.read(WMTSCapabilities);
      console.log(`## in getWmtsBaseLayers(${url} : WMTSCapabilitiesParsed : \n`, WMTSCapabilitiesParsed);
      
      const WMTSCapabilitiesSWISSTOPO = await getWMTSCapabilitiesFromUrl(urlSWISSTOPO);
      const WMTSCapabilitiesParsedSWISSTOPO = parser.read(WMTSCapabilitiesSWISSTOPO);
      console.log(`## in getWmtsBaseLayers(${urlSWISSTOPO} : WMTSCapabilitiesParsedSWISSTOPO : \n`, WMTSCapabilitiesParsedSWISSTOPO);

      arrWmtsLayers.push(
        {
          title: 'Lidar 2024',
          type: 'base',
          layer: 'orthophotos_ortho_lidar_2024',
          tile: createBaseOlLayerTile(
                  WMTSCapabilitiesParsed,
                  'Orthophoto 2024',
                  'orthophotos_ortho_lidar_2024',
                  (defaultBaseLayer === 'orthophotos_ortho_lidar_2024'),
                ),
          textStyle: {
            fill: {
              color: [255, 255, 255],
            },
            stroke: {
              color: [0, 0, 0],
              width: 0.5
            }
          }
        }
      );

      arrWmtsLayers.push(
        {
          title: 'Lidar 2016',
          type: 'base',
          layer: 'orthophotos_ortho_lidar_2016',
          tile: createBaseOlLayerTile(
                  WMTSCapabilitiesParsed,
                  'Orthophoto 2016',
                  'orthophotos_ortho_lidar_2016',
                  (defaultBaseLayer === 'orthophotos_ortho_lidar_2016'),
                ),
          textStyle: {
            fill: {
              color: [255, 255, 255],
            },
            /*
            stroke: {
              color: [0, 0, 0],
              width: 0.5
            }
            */
          }
        }
      );

      arrWmtsLayers.push(
        {
          title: 'Lidar 2012',
          type: 'base',
          layer: 'orthophotos_ortho_lidar_2012',
          tile: createBaseOlLayerTile(
                  WMTSCapabilitiesParsed,
                  'Orthophoto 2012',
                  'orthophotos_ortho_lidar_2012',
                  (defaultBaseLayer === 'orthophotos_ortho_lidar_2012'),
                ),
          textStyle: {
            fill: {
              color: [255, 255, 255],
            },
            stroke: {
              color: [0, 0, 0],
              width: 0.5
            }
          }
        }
      );

      arrWmtsLayers.push(
        {
          title: 'Plan Ville',
          type: 'base',
          layer: 'fonds_geo_osm_bdcad_couleur',
          tile: createBaseOlLayerTile(
                  WMTSCapabilitiesParsed,
                  'Plan ville',
                  'fonds_geo_osm_bdcad_couleur',
                  (defaultBaseLayer === 'fonds_geo_osm_bdcad_couleur'),
                ),
          textStyle: {
            fill: {
              color: [255, 255, 255],
            },
            stroke: {
              color: [0, 0, 0],
              width: 0.5
            }
          }
        }
      );

      arrWmtsLayers.push(
        {
          title: 'Fond cadastral',
          type: 'base',
          layer: 'fonds_geo_osm_bdcad_gris',
          tile: createBaseOlLayerTile(
                  WMTSCapabilitiesParsed,
                  'Fond cadastral',
                  'fonds_geo_osm_bdcad_gris',
                  (defaultBaseLayer === 'fonds_geo_osm_bdcad_gris'),
                ),
          textStyle: {
            fill: {
              color: [0, 0, 0],
            },
            /*
            stroke: {
              color: [255, 255, 255],
              width: 0.5
            }
            */
          }
        }
      );

      arrWmtsLayers.push(
        {
          title: 'SwissImage 2020 10cm (SWISSTOPO)',
          type: 'base',
          layer: 'ch.swisstopo.swissimage',
          tile: createBaseOlLayerTile(
                  WMTSCapabilitiesParsedSWISSTOPO,
                  'SwissImage 2020',
                  'ch.swisstopo.swissimage',
                  (defaultBaseLayer === 'ch.swisstopo.swissimage'),
                ),
          textStyle: {
            fill: {
              color: [255, 255, 255],
            },
            /*
            stroke: {
              color: [255, 255, 255],
              width: 0.5
            }
            */
          }
        }
      );

      return arrWmtsLayers;

  } catch (err) {
      const message = `###!### ERROR in getWmtsBaseLayers occured with url:${url}: error is: ${err}`;
      console.warn(message);
      return [];
  }
}

export async function createLausanneMap (
  divOfMap = 'map',
  centerOfMap = lausanneGareinMN95,
  zoomLevel = 5,
  baseLayer = defaultBaseLayer
)
{
		const arrBaseLayers = await getWmtsBaseLayers(urlLausanneMN95);
		if ((arrBaseLayers === null) || (arrBaseLayers.length < 1)) {
				console.warn('arrBaseLayers cannot be null or empty to be able to see a nice map !');
		}
    const layers = [];
    const textStyles = {};
    arrBaseLayers.forEach((layer) => {
      layers.push(layer.tile);
      textStyles[layer.layer] = layer.textStyle;
    });

		const map = new OlMap({
				target: divOfMap,
				layers: layers,
				view: new OlView({
						projection: swissProjection,
						center: centerOfMap,
						zoom: zoomLevel,
            constrainRotation: true,
				}),
		});
		setBaseLayer(map, baseLayer);

		return {map: map, layers: layers, textStyles: textStyles};
}
