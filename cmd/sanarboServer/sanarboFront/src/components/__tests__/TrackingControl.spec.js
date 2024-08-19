import { mount, config } from '@vue/test-utils';
import { describe, it, expect, beforeEach, vi } from 'vitest';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import TrackingControl from '../TrackingControl.vue';

const vuetify = createVuetify({
  components,
  directives,
})
config.global.plugins = [vuetify]

global.ResizeObserver = require('resize-observer-polyfill')

// Mock Geolocation class from ol
let positionChangeCallback;
vi.mock('ol/Geolocation.js', () => {
  return {
    default: vi.fn().mockImplementation(() => ({
      setTracking: vi.fn(),
      on: vi.fn((event, callback) => {
        if (event === 'change:position') {
          positionChangeCallback = callback;
        }
      }),
      getPosition: vi.fn(() => [0, 0]),
    })),
  };
});

describe('TrackingControl.vue', () => {
  let wrapper;
  const projection = {}; // Remplacez par un objet de projection valide si nécessaire

  beforeEach(() => {
    wrapper = mount(TrackingControl, {
      props: {
        trackingEnabled: false,
        projection: projection,
      },
    });
  });


  it('renders correctly with initial props', () => {
    expect(wrapper.exists()).toBe(true);
    expect(wrapper.find('.btn-tracking-off').exists()).toBe(true);
    expect(wrapper.find('.btn-tracking-on').exists()).toBe(false);
  });

  it('toggles tracking on button click', async () => {
    const button = wrapper.find('button');
    await button.trigger('click');

    expect(wrapper.emitted()['toggle-tracking'][0]).toEqual([true]);
    expect(wrapper.find('.btn-tracking-on').exists()).toBe(true);
    expect(wrapper.find('.btn-tracking-off').exists()).toBe(false);

    await button.trigger('click');

    expect(wrapper.emitted()['toggle-tracking'][1]).toEqual([false]);
    expect(wrapper.find('.btn-tracking-on').exists()).toBe(false);
    expect(wrapper.find('.btn-tracking-off').exists()).toBe(true);
  });

  it('emits position-changed event on geolocation position change', async () => {
    // Simulate position change
    positionChangeCallback();

    console.log(JSON.stringify(wrapper.emitted(), null, 2));
    console.dir(wrapper.emitted(), { depth: null, colors: true });
    expect(wrapper.emitted()['position-changed'][0]).toEqual([{ coords: [0, 0], zoom: 10 }]);
  });
});