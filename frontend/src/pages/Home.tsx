import React, { useEffect, useRef, useState } from 'react';
import maplibregl from 'maplibre-gl';
import { Protocol } from 'pmtiles';

import 'maplibre-gl/dist/maplibre-gl.css';
import MainButtons from './map_pages/MainButtons';

const Home: React.FC = () => {
  const mapRef = useRef<HTMLDivElement | null>(null);
  type MapMode = 'idle' | 'create-telecom' | 'create-magistral';
  const [mode, setMode] = useState<MapMode>('idle');
  const modeRef = useRef<MapMode>('idle');

  useEffect(() => {
    modeRef.current = mode;
  }, [mode]);

  useEffect(() => {
    if (!mapRef.current) return;

    const protocol = new Protocol();
    maplibregl.addProtocol('pmtiles', protocol.tile);

    const map = new maplibregl.Map({
      container: mapRef.current,
      style: '/maps/style.json',
      center: [59.6, 39.0],
      zoom: 5,
    });

    // 🔥 ШАГ 3 — ВОТ ОН
    map.on('load', async () => {
      // const icons = ['hospital', 'bus', 'town-hall', 'pharmacy', 'bank'];
      const icons = [
        'aerialway',
        'airfield',
        'airport',
        'alcohol-shop',
        'american-football',
        'amusement-park',
        'animal-shelter',
        'aquarium',
        'arrow',
        'art-gallery',
        'attraction',
        'bakery',
        'bank',
        'bank-JP',
        'bar',
        'barrier',
        'baseball',
        'basketball',
        'bbq',
        'beach',
        'beer',
        'bicycle',
        'bicycle-share',
        'blood-bank',
        'bowling-alley',
        'bridge',
        'building',
        'building-alt1',
        'bus',
        'cafe',
        'campsite',
        'car',
        'car-rental',
        'car-repair',
        'casino',
        'castle',
        'castle-JP',
        'caution',
        'cemetery',
        'cemetery-JP',
        'charging-station',
        'cinema',
        'circle',
        'circle-stroked',
        'city',
        'clothing-store',
        'college',
        'college-JP',
        'commercial',
        'communications-tower',
        'confectionery',
        'construction',
        'convenience',
        'cricket',
        'cross',
        'dam',
        'danger',
        'defibrillator',
        'dentist',
        'diamond',
        'doctor',
        'dog-park',
        'drinking-water',
        'elevator',
        'embassy',
        'emergency-phone',
        'entrance',
        'entrance-alt1',
        'farm',
        'fast-food',
        'fence',
        'ferry',
        'ferry-JP',
        'fire-station',
        'fire-station-JP',
        'fitness-centre',
        'florist',
        'fuel',
        'furniture',
        'gaming',
        'garden',
        'garden-centre',
        'gate',
        'gift',
        'globe',
        'golf',
        'grocery',
        'hairdresser',
        'hardware',
        'heart',
        'heliport',
        'highway-rest-area',
        'historic',
        'home',
        'horse-riding',
        'hospital',
        'hospital-JP',
        'hot-spring',
        'ice-cream',
        'industry',
        'information',
        'jewelry-store',
        'karaoke',
        'landmark',
        'landmark-JP',
        'landuse',
        'library',
        'lift-gate',
        'lighthouse',
        'lighthouse-JP',
        'lodging',
        'logging',
        'marae',
        'marker',
        'marker-stroked',
        'mobile-phone',
        'monument',
        'monument-JP',
        'mountain',
        'museum',
        'music',
        'natural',
        'nightclub',
        'observation-tower',
        'optician',
        'paint',
        'park',
        'park-alt1',
        'parking',
        'parking-paid',
        'pharmacy',
        'picnic-site',
        'pitch',
        'place-of-worship',
        'playground',
        'police',
        'police-JP',
        'post',
        'post-JP',
        'prison',
        'racetrack',
        'racetrack-boat',
        'racetrack-cycling',
        'racetrack-horse',
        'rail',
        'rail-light',
        'rail-metro',
        'recycling',
        'religious-buddhist',
        'religious-christian',
        'religious-jewish',
        'religious-muslim',
        'religious-shinto',
        'residential-community',
        'restaurant',
        'restaurant-bbq',
        'restaurant-noodle',
        'restaurant-pizza',
        'restaurant-seafood',
        'restaurant-sushi',
        'road-accident',
        'roadblock',
        'rocket',
        'school',
        'school-JP',
        'scooter',
        'shelter',
        'shoe',
        'shop',
        'skateboard',
        'skiing',
        'slaughterhouse',
        'slipway',
        'snowmobile',
        'soccer',
        'square',
        'square-stroked',
        'stadium',
        'star',
        'star-stroked',
        'suitcase',
        'swimming',
        'table-tennis',
        'taxi',
        'teahouse',
        'telephone',
        'tennis',
        'terminal',
        'theatre',
        'toilet',
        'toll',
        'town',
        'town-hall',
        'triangle',
        'triangle-stroked',
        'tunnel',
        'veterinary',
        'viewpoint',
        'village',
        'volleyball',
        'warehouse',
        'waste-basket',
        'watch',
        'water',
        'waterfall',
        'watermill',
        'wetland',
        'wheelchair',
        'windmill',
        'zoo',
      ];

      for (const name of icons) {
        try {
          const image = await map.loadImage(`/maps/icons/${name}.png`);

          if (!map.hasImage(name)) {
            map.addImage(name, image.data, { sdf: true });
          }
        } catch (e) {
          console.error('❌ icon load error:', name, e);
        }
      }

      // ✅ 2. Добавляем source
      map.addSource('telecom', {
        type: 'geojson',
        data: '/data/telecom_objects.geojson',
      });

      // ✅ 3. Добавляем слой ПОСЛЕ иконок
      map.addLayer({
        id: 'telecom-icons',
        type: 'symbol',
        source: 'telecom',
        minzoom: 14,
        layout: {
          'icon-image': [
            'match',
            ['get', 'type'],

            'cabinet',
            'cabinet',
            'splice',
            'splice',
            'home',
            'home',

            'marker', // fallback
          ],
          'icon-size': 1,
          'icon-allow-overlap': true,
        },
      });
    });

    // map.on('click', (e) => {
    //   const { lng, lat } = e.lngLat;
    //   console.log('📍 Координаты:', lng, lat);
    // });

    map.on('click', (e) => {
      const { lng, lat } = e.lngLat;

      if (modeRef.current === 'create-telecom') {
        console.log('📡 CREATE TELECOM:', lng, lat);

        // ⏳ дальше тут будет:
        // add point to geojson
        // or POST /nodes

        setMode('idle');
        return;
      }

      console.log('📍 Обычный клик:', lng, lat);
    });

    (window as any).map = map;

    return () => {
      map.remove();
      maplibregl.removeProtocol('pmtiles');
    };
  }, []);

  const zoomToDashoguz = () => {
    if (!(window as any).map) return;

    (window as any).map.fitBounds(
      [
        [58.5, 41.1], // юго-запад
        [60.4, 42.5], // северо-восток
      ],
      {
        padding: {
          top: 120, // 👈 под кнопку
          bottom: 40, // 👈 показываем низ
          left: 30,
          right: 30,
        },
        duration: 1200,
      },
    );
  };
  return (
    <>
      <div ref={mapRef} style={{ width: '100%', height: '100vh' }} />

      <div>
        <MainButtons zoomToDashoguz={zoomToDashoguz} mode={mode} setMode={setMode} />
      </div>
    </>
  );
};

export default Home;
