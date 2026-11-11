import { mount, flushPromises, config } from '@vue/test-utils';
import { describe, it, expect, vi } from 'vitest';
import { ref } from 'vue';
import { createVuetify } from 'vuetify';
import * as components from 'vuetify/components';
import * as directives from 'vuetify/directives';
import TreeForm from '../TreeForm.vue';
import { useFetch } from '../../composables/FetchData.js';

const vuetify = createVuetify({
  components,
  directives,
})
config.global.plugins = [vuetify]

global.ResizeObserver = require('resize-observer-polyfill')

const storedTree = {
  id: 42,
  external_id: 1042,
  name: 'Platane 42',
  description: 'Avenue de la Gare',
  is_active: true,
  is_validated: true,
  create_time: '2024-01-01T00:00:00Z',
  creator: 1007,
  comment: 'arbre remarquable',
  inactivation_reason: 'ancienne raison',
  geom: 'POINT(2538000 1152000)',
  tree_attributes: { idvalidation: 1, circonference: 120 },
};

vi.mock('../../composables/FetchData.js', () => ({
  useFetch: vi.fn(),
}));

describe('TreeForm', () => {
  it('sends back the fields it does not edit when saving', async () => {
    useFetch.mockImplementation(async (url, options) => {
      if (options.method === 'PUT') {
        return { data: ref(JSON.parse(options.body)) };
      }
      return { data: ref(structuredClone(storedTree)) };
    });
    const dictionaries = new Proxy({ validation: { data: [{ id: 1, value: 'à valider' }] } }, {
      get: (target, name) => target[name] ?? { data: [] },
    });

    const wrapper = mount(TreeForm, { props: { showForm: true, treeId: 42, dictionaries } });
    await flushPromises();
    await wrapper.vm.submitForm();

    const putCall = useFetch.mock.calls.find(([, options]) => options.method === 'PUT');
    expect(putCall, 'the form must send a PUT request').toBeDefined();
    const body = JSON.parse(putCall[1].body);
    expect(body.comment).toBe('arbre remarquable');
    expect(body.inactivation_reason).toBe('ancienne raison');
    expect(body.is_validated).toBe(false);
  });
});
