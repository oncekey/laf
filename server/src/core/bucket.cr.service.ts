import { Injectable, Logger } from '@nestjs/common'
import { GetApplicationNamespaceById } from '../common/getter'
import { KubernetesService } from './kubernetes.service'
import { CreateBucketDto } from '../buckets/dto/create-bucket.dto'
import { UpdateBucketDto } from '../buckets/dto/update-bucket.dto'
import { Bucket, BucketList } from './api/bucket.cr'

@Injectable()
export class BucketCoreService {
  logger: Logger = new Logger(BucketCoreService.name)
  constructor(private readonly k8sClient: KubernetesService) {}

  /**
   * Create a new bucket
   * @param appid
   * @param dto
   * @returns
   */
  async create(appid: string, dto: CreateBucketDto) {
    const namespace = GetApplicationNamespaceById(appid)
    const bucket = new Bucket(dto.fullname(appid), namespace)
    bucket.spec = {
      policy: dto.policy,
      storage: dto.storage,
    }

    try {
      const res = await this.k8sClient.objectApi.create(bucket)
      return Bucket.fromObject(res.body)
    } catch (error) {
      this.logger.error(error)
      return null
    }
  }

  /**
   * Query buckets of a app
   * @param appid
   * @param labelSelector
   * @returns
   */
  async findAll(appid: string, labelSelector?: string) {
    const namespace = GetApplicationNamespaceById(appid)
    const res = await this.k8sClient.customObjectApi.listNamespacedCustomObject(
      Bucket.GVK.group,
      Bucket.GVK.version,
      namespace,
      Bucket.GVK.plural,
      undefined,
      undefined,
      undefined,
      labelSelector,
    )

    return BucketList.fromObject(res.body as any)
  }

  async findOne(appid: string, name: string): Promise<Bucket> {
    const namespace = GetApplicationNamespaceById(appid)
    try {
      const res =
        await this.k8sClient.customObjectApi.getNamespacedCustomObject(
          Bucket.GVK.group,
          Bucket.GVK.version,
          namespace,
          Bucket.GVK.plural,
          name,
        )
      return Bucket.fromObject(res.body)
    } catch (err) {
      if (err?.response?.body?.reason === 'NotFound') {
        return null
      }
      this.logger.error(err, err.response?.body)
      throw err
    }
  }

  async update(bucket: Bucket, dto: UpdateBucketDto) {
    if (dto.policy) {
      bucket.spec.policy = dto.policy
    }

    if (dto.storage) {
      bucket.spec.storage = dto.storage
    }

    try {
      const res = await this.k8sClient.patchCustomObject(bucket)
      return Bucket.fromObject(res)
    } catch (error) {
      this.logger.error(error, error.response?.body)
      return null
    }
  }

  async remove(bucket: Bucket) {
    try {
      const res = await this.k8sClient.deleteCustomObject(bucket)
      return Bucket.fromObject(res)
    } catch (error) {
      this.logger.error(error, error.response?.body)
      return null
    }
  }
}
