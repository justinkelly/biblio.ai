# Copyright 2010-2019 Amazon.com, Inc. or its affiliates. All Rights Reserved.
#
# This file is licensed under the Apache License, Version 2.0 (the "License").
# You may not use this file except in compliance with the License. A copy of the
# License is located at
#
# http://aws.amazon.com/apache2.0/
#
# This file is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS
# OF ANY KIND, either express or implied. See the License for the specific
# language governing permissions and limitations under the License.

from PIL import Image
import requests
from io import BytesIO

import boto3
storingVariable = boto3.setup_default_session(region_name='us-east-1')

if __name__ == "__main__":

    # Change the value of bucket to the S3 bucket that contains your image file.
    # Change the value of photo to your image file name.

    bucket='swin-alexa'
   # photo='eureka-small.png'
    photo='george_swinburne.png'

    client=boto3.client('rekognition')


    response=client.detect_text(Image={'S3Object':{'Bucket':bucket,'Name':photo}})



    textDetections=response['TextDetections']

    print ('Detected text')
    for text in textDetections:
            print ('Detected text:' + text['DetectedText'])
            print ('Confidence: ' + "{:.2f}".format(text['Confidence']) + "%")
            print ('Id: {}'.format(text['Id']))
            if 'ParentId' in text:
                print ('Parent Id: {}'.format(text['ParentId']))
            print ('Type:' + text['Type'])
            print

