# perceptual-hash
This Library calculates the similarity between two images and returns a value in similarity in percentage

# Algorithm

- Start
- Open Images to compare
- Decode them using a suitable Decoder(`jpeg`, `png`)
- Convert them to Grayscale
- Downsize them to a 9x9(faster, less accurate) or 17x17(slightly slower, but more accurate)
- Calculate the `dHash` or `differenceHash` of each image
- Compare the two hashes and then calculate the difference in percentage b/w the two
- End

# Notes

You can provide hash length to whatever value you want. Keep in mind, that increasing hash length will decrease performance and might increase accuracy.

## Example

Say, You provided 512 bits as hash length, Then,

    512 / 2 = 256
    sqrt(256) = 16

Now the image is downsized to 17x17.

# [Inspiration](http://tech.jetsetter.com/2017/03/21/duplicate-image-detection/)

# License

MIT