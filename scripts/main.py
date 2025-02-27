import numpy as np
import cv2
import PIL, PIL.Image

def img_to_numpy_array(img):
    """
    Converts a pillow image to a numpy array representation.
    """
    # im2arr.shape: height x width x channel
    return np.array(img)[:, :, :3]  # slice to remove alpha


def img_from_numpy_array(img_arr):
    """
    Converts a numpy image representation to a pillow image.
    """
    return PIL.Image.fromarray(img_arr)

def gaussian_matrix_implementation(x_center, y_center, sigma_x, sigma_y, x_size, y_size, theta=0):
    # https://stackoverflow.com/a/67082773

    # x_center and y_center will be the center of the gaussian, theta will be the rotation angle
    # sigma_x and sigma_y will be the stdevs in the x and y axis before rotation
    # x_size and y_size give the size of the frame

    theta = 2 * np.pi * theta / 360
    x = np.arange(0, x_size, 1, float)
    y = np.arange(0, y_size, 1, float)
    y = y[:, np.newaxis]
    sx = sigma_x
    sy = sigma_y
    x0 = x_center
    y0 = y_center

    # rotation
    a = np.cos(theta) * x - np.sin(theta) * y
    b = np.sin(theta) * x + np.cos(theta) * y
    a0 = np.cos(theta) * x0 - np.sin(theta) * y0
    b0 = np.sin(theta) * x0 + np.cos(theta) * y0

    return np.exp(-(((a - a0) ** 2) / (2 * (sx ** 2)) + ((b - b0) ** 2) / (2 * (sy ** 2))))

def make_gaussian_matrix(size_x, size_y, strength=1):
    """
    Make gaussian numpy matrix given x and y dimension.
    Strength parameter specifies how sharp the gaussian function will be.
    The higher, the less sharp.
    """
    center_x = (size_x - 1) / 2
    center_y = (size_y - 1) / 2

    sigma_x = center_x * strength
    sigma_y = center_y * strength
    sigma_x = max(sigma_x, 1)
    sigma_y = max(sigma_y, 1)

    return gaussian_matrix_implementation(x_center=center_x, y_center=center_y,
                                          sigma_x=sigma_x, sigma_y=sigma_y, x_size=size_x, y_size=size_y)


def kuwahara(orig_img, method='mean', radius=3, sigma=None, grayconv=cv2.COLOR_RGB2GRAY, image_2d=None):
    """
    Get the image filtered with Kuwahara method.

    :param orig_img: original numpy image (support multichannel)
    :type orig_img: :class:`numpy.ndarray`
    :param method: method used to compute the pixels values
    :type method: "gaussian" | "mean"
    :param radius: the window radius (`winsize = 2 * radius + 1`)
    :type radius: `int`
    :param sigma: the sigma used if metod is "gaussian", automatically computed by OpenCV when `None`
    :type sigma: `float` or `None`
    :param grayconv: The OpenCV conversion code to extract grayscale image from `orig_img` (default `COLOR_BGR2GRAY`)
    :type grayconv: `int`
    :param image_2d: The 1-channel image used to compute the variance, if provided instead of `grayconv`
    :type image_2d: `numpy.ndarray`
    :returns: the filtered image
    :rtype: `numpy.ndarray`
    """

    # Filter description:
    # The Kuwahara filter works on a window divided into 4 overlapping
    # subwindows (for a 5x5 pixels example, see below). In each subwindow, the mean and
    # variance are computed. The output value (located at the center of the
    # window) is set to the mean of the subwindow with the smallest variance.
    #
    #     ( a  a  ab   b  b)
    #     ( a  a  ab   b  b)
    #     (ac ac abcd bd bd)
    #     ( c  c  cd   d  d)
    #     ( c  c  cd   d  d)
    #
    # Inspired by:
    #       python-kuwahara by Andrew Dussault, 2015 : https://github.com/adussault/python-kuwahara
    #       which it itself based on original Matlab code from Luca Balbi, 2007

    if orig_img.ndim != 2 and orig_img.ndim != 3:
        raise TypeError("Incorrect number of dimensions (excepted 2 or 3)")

    if not isinstance(radius, int):
        raise TypeError('`radius` must be int')

    if radius < 1:
        raise ValueError('`radius` must be greater or equal 1')

    if method not in ('mean', 'gaussian'):
        raise NotImplementedError('unsupported method %s' % method)

    if method == 'gaussian' and sigma is None:
        sigma = -1
        # then computed by OpenCV as : 0.3 * ((ksize - 1) * 0.5 - 1) + 0.8
        # see: https://docs.opencv.org/master/d4/d86/group__imgproc__filter.html#gac05a120c1ae92a6060dd0db190a61afa

    # convert to float32 if necessary for further math computation
    image = orig_img.astype(np.float32, copy=False)

    if image_2d is not None:
        image_2d = image_2d.astype(image.dtype, copy=False)

    # preallocate these arrays
    avgs = np.empty((4, *image.shape), dtype=image.dtype)
    print(avgs.shape)
    stddevs = np.empty((4, *image.shape[:2]), dtype=image.dtype)

    if image.ndim == 3:
        if image_2d is None:
            # NOTE this doesn't support float64
            image_2d = cv2.cvtColor(orig_img, grayconv).astype(image.dtype, copy=False)
        avgs_2d = np.empty((4, *image.shape[:2]), dtype=image.dtype)
    elif image.ndim == 2:
        image_2d = image
        avgs_2d = avgs

    print(*image_2d)
    print(avgs_2d.shape)
    # Create a pixel-by-pixel square of the image
    squared_img = image_2d ** 2

    if method == 'mean':
        kxy = np.ones(radius + 1, dtype=image.dtype) / (radius + 1)    # kernelX and kernelY (same)
    elif method == 'gaussian':
        kxy = cv2.getGaussianKernel(2 * radius + 1, sigma, ktype=cv2.CV_32F)
        kxy /= kxy[radius:].sum()   # normalize the semi-kernels
        klr = np.array([kxy[:radius+1], kxy[radius:]])
        kindexes = [[1, 1], [1, 0], [0, 1], [0, 0]]

    # the pixel position for all kernel quadrants
    shift = [(0, 0), (0,  radius), (radius, 0), (radius, radius)]

    # Calculation of averages and variances on subwindows
    for k in range(4):
        if method == 'mean':
            kx = ky = kxy
        elif method == 'gaussian':
            kx, ky = klr[kindexes[k]]
        cv2.sepFilter2D(image, -1, kx, ky, avgs[k], shift[k])
        if image_2d is not image:  # else, this is already done...
            cv2.sepFilter2D(image_2d, -1, kx, ky, avgs_2d[k], shift[k])
        cv2.sepFilter2D(squared_img, -1, kx, ky, stddevs[k], shift[k])
        stddevs[k] = stddevs[k] - avgs_2d[k] ** 2    # compute the final variance on subwindow
        #stddevs[k] = - avgs_2d[k]
    # Choice of index with minimum variance
    indices = np.argmin(stddevs, axis=0)

    # Building the filtered image
    if image.ndim == 2:
        filtered = np.take_along_axis(avgs, indices[None,...], 0).reshape(image.shape)
    else:   # then avgs.ndim == 4
        filtered = np.take_along_axis(avgs, indices[None,...,None], 0).reshape(image.shape)

    return filtered.astype(orig_img.dtype)


def apply_gaussian_kuwahara_filter(img_arr, radius=2, sigma=1.5):
    # NOTE: with sigma >= radius, this is equivalent to using 'mean' method
    # NOTE: with sigma << radius, the radius has no effect
    return kuwahara(img_arr, method='gaussian', radius=radius, sigma=sigma)


#Y←0.299⋅R+0.587⋅G+0.114⋅B



with PIL.Image.open(f'../images/test_set_pixelarts_grided/GRIDED_1.5_4.5_carps.png') as img:
    img.load()
img = img.convert()
img_arr = img_to_numpy_array(img)

# TESTING_SCAN_THROUGH_DIRECTORY("Images/test_set_pixelarts_clean")

kuwaharad_arr = apply_gaussian_kuwahara_filter(img_arr)

kuwahard_img = img_from_numpy_array(kuwaharad_arr)
kuwahard_img.save("../images/test/RESULT.png")

