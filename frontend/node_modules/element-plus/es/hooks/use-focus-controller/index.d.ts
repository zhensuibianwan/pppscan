import type { ShallowRef } from 'vue';
interface UseFocusControllerOptions {
    afterFocus?: () => void;
    afterBlur?: () => void;
}
export declare function useFocusController<T extends HTMLElement>(target: ShallowRef<T | undefined>, { afterFocus, afterBlur }?: UseFocusControllerOptions): {
    wrapperRef: ShallowRef<HTMLElement | undefined>;
    isFocused: import("vue").Ref<boolean>;
    handleFocus: (event: FocusEvent) => void;
    handleBlur: (event: FocusEvent) => void;
};
export {};
