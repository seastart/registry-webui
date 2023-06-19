/**
 * A collection of images, each of which may have multiple tags.
 */
export interface Repo {
    /** Name */
    name: string;
    /** Description */
    desc: string;
    /** Last update time (UNIX timestamp) */
    last_update: number;
    /** Tags (array) */
    tags: Tag[];
}

/**
 * Tag may have multiple images for different os/arch.
 */
export interface Tag {
    /** Name */
    name: string;
    /** Created time (UNIX timestamp) */
    created: number;
    /** Change log */
    change_log: string;
    /** Images (array) */
    images: Image[];
}

/**
 * OS/arch specified image may have multiple layers.
 */
export interface Image {
    /** Hash value */
    digest: string;
    /** CPU architecture */
    arch: string;
    /** Operating system */
    os: string;
    /** Size in bytes */
    size: number;
    /** Layers (array) */
    layers: Layer[];
}

/**
 * A layer.
 */
export interface Layer {
    /** Script content */
    script: string;
    /** Size in bytes */
    size: number;
}
