import {reactive, toRefs} from "vue";

export const useFetch = async (url, options) => {
    const state = reactive(
        {
            isLoading: true,
            hasError: false,
            errorMessage: '',
            data: [],
        }
    )



    const fetchData = async () => {
        state.isLoading = true;

        try {
            const res = await fetch(url, options);

            if (!res.ok) {
                throw new Error(res.statusText);
            }

            state.data = await res.json();

        } catch (err) {
            const typedError = err;
            state.hasError = true;
            state.errorMessage = typedError.message;
        } finally {
            state.isLoading = false;
        }
    };

        await fetchData();

        return {...toRefs(state)};
    }
